using System.Reflection;

namespace Template.Common;

public class SourceAndVersion {
    public static string GetSourceName() => Assembly.GetEntryAssembly()?.GetName().Name ?? "unknown";
    public static string GetVersion() => Assembly.GetEntryAssembly()?.GetName().Version?.ToString() ?? string.Empty;
}