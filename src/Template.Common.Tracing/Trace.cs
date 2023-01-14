using System.Diagnostics;

namespace Template.Common.Tracing;

public static class Trace {
    public static readonly ActivitySource Source =
        new(
            SourceAndVersion.SourceName,
            SourceAndVersion.Version
        );
}