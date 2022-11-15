using System.Reflection;

namespace Template.Common;

public static class SourceAndVersion {
    public static string SourceName { get; }
    public static string Version { get; }
    
    static SourceAndVersion()
    {
        var entryAssembly = Assembly.GetEntryAssembly();

        SourceName = "unknown";
        Version = string.Empty;

        if (entryAssembly is not null)
        {
            var assemblyName = entryAssembly.GetName();

            if (!string.IsNullOrEmpty(assemblyName.Name))
            {
                SourceName = assemblyName.Name;
            }

            if (assemblyName.Version is not null)
            {
                SourceName = assemblyName.Version.ToString();
            }
        }
    }
}