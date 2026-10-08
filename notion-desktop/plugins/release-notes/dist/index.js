var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
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
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  // If the importer is in node compatibility mode or this is not an ESM
  // file that has been converted to a CommonJS file using a Babel-
  // compatible transform (i.e. "__esModule" has not been set), then set
  // "default" to the CommonJS "module.exports" for node compatibility.
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/index.ts
var src_exports = {};
__export(src_exports, {
  cmpVer: () => cmpVer,
  default: () => src_default,
  notesFor: () => notesFor,
  parseNotes: () => parseNotes,
  render: () => render
});
module.exports = __toCommonJS(src_exports);
var fs = __toESM(require("fs"));
var HEAD = /^##\s+Neu in\s+(\d+(?:\.\d+)*)\s*$/;
function cmpVer(a, b) {
  const x = a.split(".").map(Number), y = b.split(".").map(Number);
  for (let i = 0; i < Math.max(x.length, y.length); i++) {
    const d = (x[i] || 0) - (y[i] || 0);
    if (d) return d;
  }
  return 0;
}
function parseNotes(md) {
  const out = [];
  let cur = null;
  for (const raw of md.replace(/\r/g, "").split("\n")) {
    const line = raw.trim();
    const m = HEAD.exec(line);
    if (m) {
      cur = { version: m[1], items: [] };
      out.push(cur);
      continue;
    }
    if (/^#{1,2}\s/.test(line)) {
      cur = null;
      continue;
    }
    if (cur && /^[-*]\s+/.test(line)) {
      const it = line.replace(/^[-*]\s+/, "").trim();
      if (it && !cur.items.includes(it)) cur.items.push(it);
    }
  }
  return out.filter((n) => n.items.length > 0);
}
function notesFor(all, version) {
  if (version) {
    const exact = all.find((n) => n.version === version);
    if (exact) return exact;
  }
  const cand = all.filter((n) => !version || cmpVer(n.version, version) <= 0).sort((a, b) => cmpVer(b.version, a.version));
  return cand[0] || all.slice().sort((a, b) => cmpVer(b.version, a.version))[0] || null;
}
function render(n) {
  if (!n) return "Keine Release-Notes gefunden.";
  return `Neu in ${n.version}:
` + n.items.map((i) => `- ${i}`).join("\n");
}
var plugin = {
  tools: [
    {
      name: "get",
      description: "Release-Notes einer AdiCode-Version (ohne version: neueste). Quelle: README-Pfad oder Markdown-Text.",
      inputSchema: {
        type: "object",
        properties: {
          version: { type: "string" },
          readme: { type: "string", description: "Pfad zur README.md" },
          markdown: { type: "string", description: "Markdown direkt (z. B. GitHub-Release-Body)" }
        }
      },
      async run(args) {
        try {
          const md = args.markdown ?? (args.readme ? fs.readFileSync(args.readme, "utf8") : "");
          if (!md) return { content: [{ type: "text", text: "readme oder markdown angeben." }], isError: true };
          return { content: [{ type: "text", text: render(notesFor(parseNotes(md), args.version)) }] };
        } catch (e) {
          return { content: [{ type: "text", text: `Fehler: ${e.message}` }], isError: true };
        }
      }
    }
  ],
  parseNotes,
  notesFor,
  render,
  cmpVer
};
var src_default = plugin;
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  cmpVer,
  notesFor,
  parseNotes,
  render
});
