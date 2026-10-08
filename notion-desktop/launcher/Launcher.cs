using System;
using System.Diagnostics;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using System.IO;
using System.IO.Compression;
using System.Management;
using System.Net;
using System.Reflection;
using System.Threading;
using System.Windows.Forms;

[assembly: AssemblyTitle("AdiCode")]
[assembly: AssemblyProduct("AdiCode (basiert auf PlazCode, GPL-3.0)")]

// Single-file launcher: embeds the full AdiCode package (app.zip), extracts it to
// %LOCALAPPDATA%\PlazCodeNotion\app, starts it and keeps watching GitHub Releases. When a
// newer release appears it replaces itself, restarts AdiCode and exits.
static class Launcher
{
    const string Repo = "369adi/PlazCode-Notion";
    const string AssetName = "PlazCode-Notion.exe";
    const string TagPrefix = "notion-desktop-v";
    const int CheckSeconds = 15;
    delegate bool EnumWindowsProc(IntPtr hwnd, IntPtr data);
    [DllImport("user32.dll")] static extern bool EnumWindows(EnumWindowsProc proc, IntPtr data);
    [DllImport("user32.dll")] static extern uint GetWindowThreadProcessId(IntPtr hwnd, out uint pid);
    [DllImport("user32.dll")] static extern bool ShowWindowAsync(IntPtr hwnd, int command);
    [DllImport("user32.dll")] static extern bool SetForegroundWindow(IntPtr hwnd);
    // Files with user state are never overwritten by an update.
    static readonly string[] Keep = { "config.json", "plazcode-settings.json" };

    static string BaseDir { get { return Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), "PlazCodeNotion"); } }
    static string Root { get { return Path.Combine(BaseDir, "app"); } }
    static string AppExe { get { return Path.Combine(Root, "PlazCode.exe"); } }
    static string Trigger { get { return Path.Combine(BaseDir, "update-now"); } }

    [STAThread]
    static int Main(string[] args)
    {
        ServicePointManager.SecurityProtocol = (SecurityProtocolType)3072; // TLS 1.2
        bool updated = Array.IndexOf(args, "--updated") >= 0;
        bool owner;
        using (Mutex mutex = new Mutex(true, "Local\\PlazCodeNotionLauncher", out owner))
        {
            if (!owner)
            {
                if (!updated)
                {
                    // A launcher is already watching; just make sure the app is visible/running.
                    if (FindApp() == null) StartApp();
                    return 0;
                }
                try { owner = mutex.WaitOne(TimeSpan.FromSeconds(90)); } catch (AbandonedMutexException) { owner = true; }
                if (!owner) return 1;
            }
            try { return Run(updated); }
            catch (Exception ex)
            {
                Log("fatal: " + ex);
                MessageBox.Show("AdiCode konnte nicht gestartet werden:\n\n" + ex.Message, "AdiCode", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return 1;
            }
            finally { try { mutex.ReleaseMutex(); } catch { } }
        }
    }

    static int Run(bool updated)
    {
        string self = Assembly.GetExecutingAssembly().Location;
        TryDelete(self + ".old");
        Version current = Normalize(Assembly.GetExecutingAssembly().GetName().Version);
        Log("launcher " + current + (updated ? " (after update)" : "") + " from " + self);

        // Kein Auto-Update mehr: installiert wird nur noch, wenn man in AdiCode auf "Aktualisieren" klickt (Datei update-now).

        EnsureExtracted(current.ToString());
        if (FindApp() == null) StartApp();

        // Neue Releases im Hintergrund vorladen, damit "Jetzt aktualisieren" sofort installiert.
        Thread pre = new Thread(delegate() { Prefetch(current); });
        pre.IsBackground = true;
        pre.Start();

        // Watch for new releases while AdiCode is running.
        while (true)
        {
            for (int i = 0; i < CheckSeconds; i++)
            {
                Thread.Sleep(1000);
                if (FindApp() == null) { Log("app closed - launcher exits"); return 0; }
                if (File.Exists(Trigger))
                {
                    // "Jetzt aktualisieren" in der App (Seite Updates).
                    TryDelete(Trigger);
                    Log("manual update check");
                    if (TryUpdate(self, current)) return 0;
                    Log("no newer release than " + current);
                }
            }
        }
    }

    // ---------- update ----------

    static bool TryUpdate(string self, Version current)
    {
        string tag; Version latest;
        if (!LatestRelease(out tag, out latest) || latest.CompareTo(current) <= 0) return false;
        Log("update available: " + current + " -> " + latest);
        string tmp = Path.Combine(Path.GetTempPath(), "PlazCode-Notion-" + latest + ".exe");
        try
        {
            // Sofort-Update: im Hintergrund vorgeladene exe nutzen statt erst jetzt 40 MB zu laden.
            string staged = StagedPath(latest);
            if (LooksLikeExe(staged)) { File.Copy(staged, tmp, true); Log("using prefetched " + staged); }
            else
            using (WebClient wc = new WebClient())
            {
                wc.Headers.Add("User-Agent", "PlazCode-Notion-Launcher");
                wc.DownloadFile("https://github.com/" + Repo + "/releases/download/" + tag + "/" + AssetName, tmp);
            }
            if (!LooksLikeExe(tmp)) { Log("download is not a valid exe"); TryDelete(tmp); return false; }

            string target = self;
            try
            {
                TryDelete(self + ".old");
                File.Move(self, self + ".old");
                try { File.Move(tmp, self); }
                catch { File.Move(self + ".old", self); throw; }
            }
            catch (Exception ex)
            {
                // e.g. read-only folder: install next to the app data instead.
                Log("in-place replace failed (" + ex.Message + "), using " + BaseDir);
                target = Path.Combine(BaseDir, AssetName);
                Directory.CreateDirectory(BaseDir);
                File.Copy(tmp, target, true);
                TryDelete(tmp);
            }
            // Task #2: alte Instanz jetzt sofort beenden, damit nichts versteckt im Tray weiterlaeuft.
            StopAppGraceful();
            ProcessStartInfo psi = new ProcessStartInfo(target, "--updated");
            psi.UseShellExecute = false;
            psi.WorkingDirectory = Path.GetDirectoryName(target);
            Process.Start(psi);
            Log("started updated launcher " + target);
            return true;
        }
        catch (Exception ex)
        {
            Log("update failed: " + ex.Message);
            TryDelete(tmp);
            return false;
        }
    }

    static string StagedPath(Version v) { return Path.Combine(BaseDir, "staged-" + v + ".exe"); }

    static void Prefetch(Version current)
    {
        while (true)
        {
            try
            {
                string tag; Version latest;
                if (LatestRelease(out tag, out latest) && latest.CompareTo(current) > 0)
                {
                    string staged = StagedPath(latest);
                    if (!LooksLikeExe(staged))
                    {
                        foreach (string old in Directory.GetFiles(BaseDir, "staged-*")) TryDelete(old);
                        string part = staged + ".part";
                        using (WebClient wc = new WebClient())
                        {
                            wc.Headers.Add("User-Agent", "PlazCode-Notion-Launcher");
                            wc.DownloadFile("https://github.com/" + Repo + "/releases/download/" + tag + "/" + AssetName, part);
                        }
                        if (LooksLikeExe(part)) { File.Move(part, staged); Log("prefetched " + latest); } else TryDelete(part);
                    }
                }
            }
            catch (Exception ex) { Log("prefetch failed: " + ex.Message); }
            Thread.Sleep(60000);
        }
    }

    // Reads the tag of the latest release from the redirect of /releases/latest (no API rate limit).
    static bool LatestRelease(out string tag, out Version version)
    {
        tag = null; version = null;
        try
        {
            HttpWebRequest req = (HttpWebRequest)WebRequest.Create("https://github.com/" + Repo + "/releases/latest");
            req.Method = "HEAD";
            req.AllowAutoRedirect = false;
            req.UserAgent = "PlazCode-Notion-Launcher";
            req.Timeout = 15000;
            using (HttpWebResponse res = (HttpWebResponse)req.GetResponse())
            {
                string loc = res.Headers["Location"];
                if (string.IsNullOrEmpty(loc)) return false;
                int i = loc.LastIndexOf("/tag/", StringComparison.Ordinal);
                if (i < 0) return false;
                tag = Uri.UnescapeDataString(loc.Substring(i + 5));
                if (!tag.StartsWith(TagPrefix, StringComparison.Ordinal)) return false;
                version = Normalize(new Version(tag.Substring(TagPrefix.Length)));
                return true;
            }
        }
        catch (Exception ex) { Log("update check failed: " + ex.Message); return false; }
    }

    static Version Normalize(Version v)
    {
        return new Version(v.Major, v.Minor, Math.Max(v.Build, 0));
    }

    static bool LooksLikeExe(string path)
    {
        FileInfo f = new FileInfo(path);
        if (!f.Exists || f.Length < 1024 * 1024) return false;
        using (FileStream s = f.OpenRead()) { return s.ReadByte() == 'M' && s.ReadByte() == 'Z'; }
    }

    // ---------- app ----------

    static void EnsureExtracted(string version)
    {
        string marker = Path.Combine(Root, ".version");
        if (File.Exists(AppExe) && File.Exists(marker) && File.ReadAllText(marker).Trim() == version) return;
        Log("installing app " + version);
        StopApp();
        Directory.CreateDirectory(Root);
        string rootFull = Path.GetFullPath(Root);
        using (Stream s = Assembly.GetExecutingAssembly().GetManifestResourceStream("app.zip"))
        using (ZipArchive zip = new ZipArchive(s, ZipArchiveMode.Read))
        {
            foreach (ZipArchiveEntry e in zip.Entries)
            {
                string dest = Path.GetFullPath(Path.Combine(rootFull, e.FullName));
                if (!dest.StartsWith(rootFull, StringComparison.OrdinalIgnoreCase)) continue;
                if (e.FullName.EndsWith("/") || e.FullName.EndsWith("\\")) { Directory.CreateDirectory(dest); continue; }
                if (File.Exists(dest) && Array.IndexOf(Keep, e.FullName) >= 0) continue;
                Directory.CreateDirectory(Path.GetDirectoryName(dest));
                e.ExtractToFile(dest, true);
            }
        }
        File.WriteAllText(marker, version);
    }

    static Process FindApp()
    {
        foreach (Process p in Process.GetProcessesByName("PlazCode"))
        {
            try { if (p.MainModule.FileName.StartsWith(Root, StringComparison.OrdinalIgnoreCase)) return p; }
            catch { }
        }
        return null;
    }

    static void StartApp()
    {
        ProcessStartInfo psi = new ProcessStartInfo(AppExe);
        psi.WorkingDirectory = Root;
        psi.UseShellExecute = false;
        Process.Start(psi);
        Log("app started");
    }

    static void RestoreCoWorkWindows()
    {
        try
        {
            HashSet<uint> pids = new HashSet<uint>();
            using (ManagementObjectSearcher q = new ManagementObjectSearcher("SELECT ProcessId,Name,CommandLine FROM Win32_Process"))
            foreach (ManagementObject o in q.Get())
            {
                string name = Convert.ToString(o["Name"] ?? "").ToLowerInvariant();
                string cmd = Convert.ToString(o["CommandLine"] ?? "").ToLowerInvariant();
                bool browser = name.Contains("chrome") || name.Contains("msedge") || name.Contains("brave");
                if (browser && cmd.Contains("plazcodenotion\\profiles\\agent-")) pids.Add(Convert.ToUInt32(o["ProcessId"]));
            }
            EnumWindows(delegate(IntPtr hwnd, IntPtr data) { uint pid; GetWindowThreadProcessId(hwnd, out pid); if (pids.Contains(pid)) { ShowWindowAsync(hwnd, 9); SetForegroundWindow(hwnd); } return true; }, IntPtr.Zero);
            Log("restored Co-Work browser windows");
        }
        catch (Exception ex) { Log("restore Co-Work windows failed: " + ex.Message); }
    }

    // Task #2: erst freundlich (CloseMainWindow, max 2,5 s), sonst hart inkl. Kindprozesse (ngrok, Add-ons); Co-Work-Browser bleiben.
    static void StopAppGraceful()
    {
        RestoreCoWorkWindows();
        Process p;
        int guard = 0;
        while ((p = FindApp()) != null && guard++ < 10)
        {
            bool asked = false;
            try { asked = p.CloseMainWindow(); } catch { }
            bool gone = false;
            try { gone = asked && p.WaitForExit(2500); } catch { }
            if (!gone) { Log("stopping app pid " + p.Id + " for update"); KillTree(p.Id); }
            Thread.Sleep(200);
        }
    }
    static void StopApp()
    {
        RestoreCoWorkWindows();
        Process p;
        while ((p = FindApp()) != null)
        {
            Log("stopping app pid " + p.Id);
            KillTree(p.Id);
            Thread.Sleep(500);
        }
    }

    // Kills a process and all its children (ngrok, MCP add-ons, ...), children first.
    static void KillTree(int pid)
    {
        try
        {
            using (ManagementObjectSearcher q = new ManagementObjectSearcher("SELECT ProcessId,Name,CommandLine FROM Win32_Process WHERE ParentProcessId=" + pid))
            {
                foreach (ManagementObject o in q.Get())
                {
                    string name = Convert.ToString(o["Name"] ?? "").ToLowerInvariant();
                    string cmd = Convert.ToString(o["CommandLine"] ?? "").ToLowerInvariant();
                    bool browser = name.Contains("chrome") || name.Contains("msedge") || name.Contains("brave");
                    bool coworkProfile = cmd.Contains("plazcodenotion\\profiles\\agent-");
                    if (browser && coworkProfile) { Log("preserving Co-Work browser pid " + o["ProcessId"]); continue; }
                    KillTree(Convert.ToInt32(o["ProcessId"]));
                }
            }
        }
        catch { }
        try { Process p = Process.GetProcessById(pid); p.Kill(); p.WaitForExit(5000); } catch { }
    }

    // ---------- helpers ----------

    static void TryDelete(string path)
    {
        try { if (File.Exists(path)) File.Delete(path); } catch { }
    }

    static void Log(string message)
    {
        try
        {
            string dir = Path.Combine(BaseDir, "logs");
            Directory.CreateDirectory(dir);
            File.AppendAllText(Path.Combine(dir, "launcher.log"), DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss") + "  " + message + Environment.NewLine);
        }
        catch { }
    }
}
