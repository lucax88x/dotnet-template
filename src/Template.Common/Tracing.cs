using System.Diagnostics;

namespace Template.Common;

public static class Tracing {
    public static readonly ActivitySource Source =
        new(
            SourceAndVersion.GetSourceName(),
            SourceAndVersion.GetVersion()
        );
}