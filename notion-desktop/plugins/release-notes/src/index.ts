// Release-Notes als Plugin (ADR work/upd134/hybrid_adr.md).
// Liest die Abschnitte "## Neu in <version>" aus dem README (lokal) bzw. aus GitHub-Release-Bodies
// und liefert saubere, deduplizierte Notes. Fallback: neueste Version mit Notes (wie notes_any in Rust).
import * as fs from "fs";

export type Notes = { version: string; items: string[] };

const HEAD = /^##\s+Neu in\s+(\d+(?:\.\d+)*)\s*$/;

export function cmpVer(a: string, b: string): number {
  const x = a.split(".").map(Number), y = b.split(".").map(Number);
  for (let i = 0; i < Math.max(x.length, y.length); i++) {
    const d = (x[i] || 0) - (y[i] || 0);
    if (d) return d;
  }
  return 0;
}

/** Alle "## Neu in X"-Abschnitte; Punkte pro Version dedupliziert (Reihenfolge bleibt). */
export function parseNotes(md: string): Notes[] {
  const out: Notes[] = [];
  let cur: Notes | null = null;
  for (const raw of md.replace(/\r/g, "").split("\n")) {
    const line = raw.trim();
    const m = HEAD.exec(line);
    if (m) { cur = { version: m[1], items: [] }; out.push(cur); continue; }
    if (/^#{1,2}\s/.test(line)) { cur = null; continue; }
    if (cur && /^[-*]\s+/.test(line)) {
      const it = line.replace(/^[-*]\s+/, "").trim();
      if (it && !cur.items.includes(it)) cur.items.push(it);
    }
  }
  return out.filter(n => n.items.length > 0);
}

/** Notes fuer genau diese Version, sonst die neueste vorhandene (<= version, falls angegeben). */
export function notesFor(all: Notes[], version?: string): Notes | null {
  if (version) {
    const exact = all.find(n => n.version === version);
    if (exact) return exact;
  }
  const cand = all.filter(n => !version || cmpVer(n.version, version) <= 0).sort((a, b) => cmpVer(b.version, a.version));
  return cand[0] || all.slice().sort((a, b) => cmpVer(b.version, a.version))[0] || null;
}

export function render(n: Notes | null): string {
  if (!n) return "Keine Release-Notes gefunden.";
  return `Neu in ${n.version}:\n` + n.items.map(i => `- ${i}`).join("\n");
}

const plugin = {
  tools: [
    {
      name: "get",
      description: "Release-Notes einer AdiCode-Version (ohne version: neueste). Quelle: README-Pfad oder Markdown-Text.",
      inputSchema: {
        type: "object",
        properties: {
          version: { type: "string" },
          readme: { type: "string", description: "Pfad zur README.md" },
          markdown: { type: "string", description: "Markdown direkt (z. B. GitHub-Release-Body)" },
        },
      },
      async run(args: { version?: string; readme?: string; markdown?: string }) {
        try {
          const md = args.markdown ?? (args.readme ? fs.readFileSync(args.readme, "utf8") : "");
          if (!md) return { content: [{ type: "text", text: "readme oder markdown angeben." }], isError: true };
          return { content: [{ type: "text", text: render(notesFor(parseNotes(md), args.version)) }] };
        } catch (e) {
          return { content: [{ type: "text", text: `Fehler: ${(e as Error).message}` }], isError: true };
        }
      },
    },
  ],
  parseNotes, notesFor, render, cmpVer,
};

export default plugin;
