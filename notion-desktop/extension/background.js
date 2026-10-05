const DEFAULT_URL = "ws://127.0.0.1:8787/extension/ws";
let socket = null;
let reconnectDelay = 1500;
let reconnectTimer = null;
let heartbeatTimer = null;
const directSessions = new Map();

async function settings() {
  return chrome.storage.local.get({bridgeUrl:DEFAULT_URL,bridgeToken:"",enabledAgents:{}});
}

async function activeTab() {
  const [tab] = await chrome.tabs.query({active:true,lastFocusedWindow:true});
  return tab?.id ? tab : null;
}

async function pageInfo(tabId) {
  try { return await chrome.tabs.sendMessage(tabId, {type:"GET_WEB_AGENT_INFO"}); }
  catch { return null; }
}

async function enabledAgents() {
  const {enabledAgents:records = {}} = await settings();
  const agents = [];
  for (const [tabId, record] of Object.entries(records)) {
    try {
      const tab = await chrome.tabs.get(Number(tabId));
      const info = await pageInfo(tab.id);
      if (!info?.supported) continue;
      agents.push({...record,...info,id:`agent-${tab.id}`,tabId:tab.id});
    } catch {}
  }
  return agents;
}

async function toggleActiveAgent() {
  const tab = await activeTab();
  if (!tab) return {error:"No active browser tab is available."};
  const info = await pageInfo(tab.id);
  if (!info?.supported) return {error:info?.reason || "Open a supported signed-in AI chat page first."};
  const {enabledAgents:records = {}} = await settings();
  const key = String(tab.id);
  if (records[key]) {
    await stopDirectSession(tab.id);
    delete records[key];
    await chrome.storage.local.set({enabledAgents:records});
    return {enabled:false,info};
  }
  records[key] = {id:`agent-${tab.id}`,site:info.site,name:info.name,enabledAt:Date.now()};
  await chrome.storage.local.set({enabledAgents:records});
  return {enabled:true,info};
}

async function stopDirectSession(tabId) {
  try { await chrome.tabs.sendMessage(tabId,{type:"DIRECT_AGENT_STOP"}); } catch {}
  const session = directSessions.get(tabId);
  if (session) {
    await directAgentApi("stop",session);
    directSessions.delete(tabId);
  }
}

async function runOnAgent(agentId, message) {
  const agents = await enabledAgents();
  const agent = agents.find(item => item.id === agentId);
  if (!agent) return {error:"This page is not enabled as an agent. Open it and choose Run agent in the extension."};
  try {
    const request = message.action === "run"
      ? {type:"WEB_AGENT_RUN",task:message.task,wait_seconds:message.wait_seconds}
      : {type:"WEB_MCP_TOOL",action:message.action,text:message.text,wait_seconds:message.wait_seconds};
    return await chrome.tabs.sendMessage(agent.tabId, request);
  } catch {
    return {error:"The enabled agent page is no longer available. Reload the page and enable it again."};
  }
}

async function directAgentApi(path, body) {
  const config = await settings();
  if (!config.bridgeToken) return {error:"Add the bridge token in the extension options (token from AdiCode > Notion AI) first. The direct Roblox agent needs the local bridge."};
  const base = (config.bridgeUrl || DEFAULT_URL)
    .replace(/^ws:/i,"http:").replace(/^wss:/i,"https:").replace(/\/extension\/ws\/?$/i,"");
  try {
    const response = await fetch(`${base}/direct-agent/${path}`, {
      method:"POST",
      headers:{"Authorization":`Bearer ${config.bridgeToken}`,"Content-Type":"application/json"},
      body:JSON.stringify(body),
    });
    const result = await response.json();
    return response.ok ? result : {error:result.error || `Bridge returned HTTP ${response.status}.`};
  } catch (error) { return {error:`Cannot reach the local Roblox bridge: ${String(error?.message || error)}`}; }
}

async function status(state, detail = "") {
  await chrome.storage.local.set({connectionState:state,connectionDetail:detail});
}

function scheduleReconnect() {
  clearTimeout(reconnectTimer);
  reconnectTimer = setTimeout(connect, reconnectDelay);
  reconnectDelay = Math.min(30000, Math.round(reconnectDelay * 1.7));
}

async function connect() {
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) return;
  const config = await settings();
  if (!config.bridgeToken) {
    await status("needs_setup", "Add the bridge token in the extension options (token from AdiCode > Notion AI).");
    return;
  }
  let ws;
  try { ws = new WebSocket(config.bridgeUrl || DEFAULT_URL); }
  catch (error) { await status("disconnected", String(error)); scheduleReconnect(); return; }
  socket = ws;
  ws.onopen = () => {
    reconnectDelay = 1500;
    ws.send(JSON.stringify({type:"hello",token:config.bridgeToken}));
    clearInterval(heartbeatTimer);
    heartbeatTimer = setInterval(() => { if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({type:"heartbeat"})); }, 15000);
  };
  ws.onmessage = async event => {
    let message;
    try { message = JSON.parse(event.data); } catch { return; }
    if (message.type === "ready") {
      await status("connected", "Connected. Enable a supported chat tab with Activer cet agent before using web_agents.");
      return;
    }
    if (message.type !== "request" || typeof message.id !== "string") return;
    try {
      const result = await handleRequest(message);
      if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({type:"response",id:message.id,result}));
    } catch (error) {
      if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({type:"response",id:message.id,result:{error:String(error?.message || error)}}));
    }
  };
  ws.onerror = () => status("disconnected", "Cannot reach the local bridge. Check that it is running.");
  ws.onclose = async event => {
    clearInterval(heartbeatTimer);
    if (socket === ws) socket = null;
    await status(event.code === 4401 ? "unauthorized" : "disconnected", event.code === 4401 ? "The bridge token was rejected. Update it in the extension options (token from AdiCode > Notion AI)." : "Waiting for the local bridge.");
    scheduleReconnect();
  };
}

async function handleRequest(message) {
  if (message.tool === "web_sites") return {error:"web_sites does not require a browser request."};
  const args = message.arguments || {};
  if (message.tool === "web_agents" && args.action === "list") return {agents:await enabledAgents()};
  if (message.tool === "web_agents") return runOnAgent(args.agent_id, args);
  const tab = await activeTab();
  if (!tab) return {error:"No active browser tab is available."};
  const {enabledAgents:records = {}} = await settings();
  const record = records[String(tab.id)];
  if (!record) return {error:"Enable this page with Run agent in the extension before using web_chat."};
  return runOnAgent(`agent-${tab.id}`, args);
}

chrome.runtime.onInstalled.addListener(connect);
chrome.runtime.onStartup.addListener(connect);
chrome.storage.onChanged.addListener((changes, area) => {
  if (area === "local" && (changes.bridgeUrl || changes.bridgeToken)) {
    if (socket) socket.close();
    else connect();
  }
});
chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (message?.type === "GET_ACTIVE_AGENT_STATE") {
    activeTab().then(async tab => {
      if (!tab) return sendResponse({error:"No active tab."});
      const info = await pageInfo(tab.id);
      const {enabledAgents:records = {}} = await settings();
      sendResponse({info,enabled:!!records[String(tab.id)]});
    });
    return true;
  }
  if (message?.type === "TOGGLE_ACTIVE_AGENT") { toggleActiveAgent().then(sendResponse); return true; }
  if (message?.type === "RUN_DIRECT_TASK") {
    activeTab().then(async tab => {
      if (!tab) return sendResponse({error:"No active tab."});
      const {enabledAgents:records = {}} = await settings();
      if (!records[String(tab.id)]) return sendResponse({error:"Enable this page with Activer cet agent before starting a direct Roblox agent."});
      try {
        const result = await chrome.tabs.sendMessage(tab.id, {type:"DIRECT_AGENT_START",task:message.task,agent_id:`agent-${tab.id}`});
        sendResponse(result);
      } catch { sendResponse({error:"Refresh the AI page, enable it as an agent, then try again."}); }
    });
    return true;
  }
  if (message?.type === "DIRECT_AGENT_API") {
    const paths = {start:"start",tool:"tool",heartbeat:"heartbeat",stop:"stop"};
    if (!Object.hasOwn(paths,message.path)) { sendResponse({error:"Unknown direct-agent request."}); return; }
    const tabId = _sender?.tab?.id;
    if (!Number.isInteger(tabId)) { sendResponse({error:"Direct Roblox tools are available only from an enabled chat tab."}); return; }
    (async () => {
      const body = message.body || {};
      const {enabledAgents:records = {}} = await settings();
      const info = await pageInfo(tabId);
      const agentId = `agent-${tabId}`;
      if (!records[String(tabId)] || !info?.supported) return {error:"Enable this supported chat tab before using the direct Roblox agent."};
      if (body.agent_id !== agentId) return {error:"This direct-agent session belongs to a different tab."};
      if (message.path === "start" && body.provider !== info.site) return {error:"The direct-agent provider does not match this chat tab."};
      const result = await directAgentApi(paths[message.path],body);
      if (!result?.error && message.path === "start") directSessions.set(tabId,{session_id:result.session_id,agent_id:agentId});
      if (message.path === "stop") directSessions.delete(tabId);
      return result;
    })().then(sendResponse).catch(error => sendResponse({error:String(error?.message || error)}));
    return true;
  }
  if (message?.type === "GET_WEB_MCP_STATUS") {
    settings().then(async value => {
      const state = await chrome.storage.local.get({connectionState:"disconnected",connectionDetail:""});
      sendResponse({...state,configured:!!value.bridgeToken});
    });
    return true;
  }
  if (message?.type === "RECONNECT_WEB_MCP") { connect(); sendResponse({ok:true}); }
});

chrome.tabs.onRemoved.addListener(async tabId => {
  const session = directSessions.get(tabId);
  if (session) { await directAgentApi("stop",session); directSessions.delete(tabId); }
  const {enabledAgents:records = {}} = await settings();
  delete records[String(tabId)];
  await chrome.storage.local.set({enabledAgents:records});
});

chrome.tabs.onUpdated.addListener(async (tabId, changeInfo) => {
  if (!changeInfo.url) return;
  const {enabledAgents:records = {}} = await settings();
  if (!records[String(tabId)]) return;
  const info = await pageInfo(tabId);
  if (!info?.supported) {
    const session = directSessions.get(tabId);
    if (session) { await directAgentApi("stop",session); directSessions.delete(tabId); }
    delete records[String(tabId)];
    await chrome.storage.local.set({enabledAgents:records});
  } else {
    records[String(tabId)] = {...records[String(tabId)],site:info.site,name:info.name};
    await chrome.storage.local.set({enabledAgents:records});
  }
});

connect();
