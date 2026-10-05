using System;
using System.Diagnostics;
using System.IO;
using System.IO.Compression;
using System.Reflection;
using System.Windows.Forms;

[assembly: AssemblyTitle("PlazCode Notion")]
[assembly: AssemblyProduct("PlazCode Notion (inoffizieller Fork)")]

static class Launcher
{
    // Dateien mit Benutzerzustand werden bei Updates nicht ueberschrieben.
    static readonly string[] Keep = { "config.json", "plazcode-settings.json" };

    [STAThread]
    static int Main(string[] args)
    {
        try
        {
            Assembly asm = Assembly.GetExecutingAssembly();
            string version = asm.GetName().Version.ToString();
            string root = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), "PlazCodeNotion", "app");
            string marker = Path.Combine(root, ".version");
            string exe = Path.Combine(root, "PlazCode.exe");

            bool current = File.Exists(exe) && File.Exists(marker) && File.ReadAllText(marker).Trim() == version;
            if (!current)
            {
                StopRunning(root);
                Directory.CreateDirectory(root);
                using (Stream s = asm.GetManifestResourceStream("app.zip"))
                using (ZipArchive zip = new ZipArchive(s, ZipArchiveMode.Read))
                {
                    foreach (ZipArchiveEntry e in zip.Entries)
                    {
                        string dest = Path.GetFullPath(Path.Combine(root, e.FullName));
                        if (!dest.StartsWith(root, StringComparison.OrdinalIgnoreCase)) continue;
                        if (e.FullName.EndsWith("/") || e.FullName.EndsWith("\\")) { Directory.CreateDirectory(dest); continue; }
                        if (File.Exists(dest) && Array.IndexOf(Keep, e.FullName) >= 0) continue;
                        Directory.CreateDirectory(Path.GetDirectoryName(dest));
                        e.ExtractToFile(dest, true);
                    }
                }
                File.WriteAllText(marker, version);
            }

            ProcessStartInfo psi = new ProcessStartInfo(exe);
            psi.WorkingDirectory = root;
            psi.UseShellExecute = false;
            psi.Arguments = string.Join(" ", Array.ConvertAll(args, a => "\"" + a + "\""));
            Process.Start(psi);
            return 0;
        }
        catch (Exception ex)
        {
            MessageBox.Show("PlazCode Notion konnte nicht gestartet werden:\n\n" + ex.Message, "PlazCode Notion", MessageBoxButtons.OK, MessageBoxIcon.Error);
            return 1;
        }
    }

    static void StopRunning(string root)
    {
        foreach (Process p in Process.GetProcessesByName("PlazCode"))
        {
            try
            {
                if (p.MainModule.FileName.StartsWith(root, StringComparison.OrdinalIgnoreCase)) { p.Kill(); p.WaitForExit(5000); }
            }
            catch { }
        }
    }
}
