const state = document.querySelector("#state");
const detail = document.querySelector("#detail");
const enable = document.querySelector("#enable");
const run = document.querySelector("#run");
const task = document.querySelector("#task");
const result = document.querySelector("#result");

async function refresh() {
  const [page, connection] = await Promise.all([
    chrome.runtime.sendMessage({type:"GET_ACTIVE_AGENT_STATE"}),
    chrome.runtime.sendMessage({type:"GET_WEB_MCP_STATUS"}),
  ]);
  if (!page?.info?.supported) {
    state.textContent = "Seite nicht unterstützt";
    detail.textContent = page?.info?.reason || connection.connectionDetail || "Öffne in diesem Tab einen unterstützten KI-Chat.";
    enable.disabled = true;
    run.disabled = true;
    return;
  }
  state.textContent = `${page.info.name}${page.enabled ? " · Agent freigegeben" : " · Agent nicht freigegeben"}`;
  detail.textContent = `PlazCode: ${connection.configured ? connection.connectionState || "getrennt" : "Token fehlt"}. Der Direkt-Modus braucht PlazCode Notion und ein verbundenes Roblox Studio; Notion ist optional.`;
  enable.textContent = page.enabled ? "Freigabe entfernen" : "Agent freigeben";
  enable.disabled = false;
  run.disabled = !page.enabled;
}

refresh();

enable.addEventListener("click", async () => {
  const response = await chrome.runtime.sendMessage({type:"TOGGLE_ACTIVE_AGENT"});
  if (response?.error) result.textContent = response.error;
  else result.textContent = response.enabled ? "Dieser Tab ist jetzt in web_agents verfügbar." : "Freigabe für diesen Tab entfernt.";
  await refresh();
});

run.addEventListener("click", async () => {
  if (!task.value.trim()) { result.textContent = "Gib zuerst eine Aufgabe ein."; return; }
  run.disabled = true;
  result.textContent = "Der KI-Chat bearbeitet die Aufgabe…";
  try {
    const response = await chrome.runtime.sendMessage({type:"RUN_DIRECT_TASK",task:task.value.trim()});
    result.textContent = response.error || response.message || response.text || "Die Seite hat keinen Text zurückgegeben.";
  } catch (error) { result.textContent = String(error); }
  await refresh();
});

document.querySelector("#options").addEventListener("click", () => chrome.runtime.openOptionsPage());