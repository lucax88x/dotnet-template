using System.Diagnostics;

namespace Template.Common.Tracing;

public static class Tracing {
    public static readonly ActivitySource Source =
        new(
            SourceAndVersion.SourceName,
            SourceAndVersion.Version
        );
}