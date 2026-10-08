var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/index.ts
var src_exports = {};
__export(src_exports, {
  KEEP: () => KEEP,
  MIN_SCROLL: () => MIN_SCROLL,
  default: () => src_default,
  turboJs: () => turboJs
});
module.exports = __toCommonJS(src_exports);
var KEEP = 6;
var MIN_SCROLL = 5e3;
function turboJs(keep = KEEP, minScroll = MIN_SCROLL) {
  return `(()=>{if(window.__adiTurbo)return;window.__adiTurbo=1;const KEEP=${keep},MIN=${minScroll};const st=document.createElement('style');st.textContent='.adi-cv{content-visibility:auto;contain-intrinsic-size:auto 600px}';let list=null;const big=e=>{let b=null,h=0;for(const c of e.children){const x=c.offsetHeight;if(x>h){h=x;b=c}}return b};const find=()=>{let best=null,h=0;for(const s of document.querySelectorAll('.notion-scroller')){if(s.scrollHeight>h){h=s.scrollHeight;best=s}}if(!best||h<MIN)return null;let n=best;for(let i=0;i<16&&n;i++){if(n.children.length>=10)return n;n=big(n)}return null};const tick=()=>{try{if(!st.isConnected)(document.head||document.documentElement).appendChild(st);if(!list||!list.isConnected||list.children.length<10)list=find();if(list){const k=list.children,n=k.length-KEEP;for(let i=0;i<k.length;i++)k[i].classList.toggle('adi-cv',i<n)}}catch(e){}};const loop=()=>{tick();setTimeout(()=>{if(window.requestIdleCallback)requestIdleCallback(loop,{timeout:3000});else loop()},5000)};setTimeout(loop,3000)})();`;
}
var plugin = {
  inject: [{ name: "chat-turbo", js: turboJs(), when: "both" }],
  tools: [
    {
      name: "status",
      description: "Chat-Turbo: zeigt pro offenem Chat-Tab, ob Turbo aktiv ist und wie viele alte Nachrichten nicht gerendert werden.",
      inputSchema: { type: "object", properties: {} },
      async run(_args, host) {
        if (!host?.cdp_eval) return { content: [{ type: "text", text: "Host ohne cdp_eval." }] };
        const r = await host.cdp_eval("all", "({turbo:!!window.__adiTurbo,skipped:document.querySelectorAll('.adi-cv').length})");
        return { content: [{ type: "text", text: JSON.stringify(r) }] };
      }
    }
  ],
  turboJs
};
var src_default = plugin;
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  KEEP,
  MIN_SCROLL,
  turboJs
});
