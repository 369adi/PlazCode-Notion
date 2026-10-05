const url = document.querySelector("#url");
const token = document.querySelector("#token");
const notion = document.querySelector("#notion");
const instructions = document.querySelector("#instructions");
const status = document.querySelector("#status");

chrome.storage.local.get({bridgeUrl:"ws://127.0.0.1:8787/extension/ws",bridgeToken:"",notionEnabled:false,agentInstructions:""}).then(config => {
  url.value = config.bridgeUrl;
  token.value = config.bridgeToken;
  notion.checked = config.notionEnabled;
  instructions.value = config.agentInstructions;
});

document.querySelector("#save").addEventListener("click", async () => {
  if (!url.value.startsWith("ws://127.0.0.1:") && !url.value.startsWith("ws://localhost:")) { status.textContent = "Aus Sicherheitsgründen muss die Adresse auf 127.0.0.1 oder localhost zeigen."; return; }
  await chrome.storage.local.set({bridgeUrl:url.value.trim(),bridgeToken:token.value.trim(),notionEnabled:notion.checked,agentInstructions:instructions.value.slice(0,8000)});
  if (token.value.trim()) {
    await chrome.runtime.sendMessage({type:"RECONNECT_WEB_MCP"});
    status.textContent = "Gespeichert. Den Verbindungsstatus siehst du im Menü der Erweiterung.";
  } else {
    status.textContent = "Anweisungen gespeichert. Trage den Erweiterungs-Token ein, um dich mit PlazCode zu verbinden.";
  }
});