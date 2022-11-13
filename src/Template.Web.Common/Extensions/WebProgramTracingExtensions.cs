using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Template.Common;
using Template.Web.Common.Configs;

namespace Template.Web.Common.Extensions;

public static class WebProgramTracingExtensions {
    public static WebApplicationBuilder AddTemplateTracing(this WebApplicationBuilder webApplicationBuilder)
    {
        var jaegerConfig = webApplicationBuilder.Configuration.GetSection(JaegerConfig.Section).Get<JaegerConfig>();

        if (jaegerConfig is not null)
        {
            webApplicationBuilder.Services
                .AddOpenTelemetryTracing(
                    builder =>
                    {
                        var sourceName = SourceAndVersion.GetSourceName();
                        var version = SourceAndVersion.GetVersion();

                        builder
                            .AddSource(sourceName)
                            .SetResourceBuilder(
                                ResourceBuilder
                                    .CreateDefault()
                                    .AddService(sourceName, serviceVersion: version)
                            )
                            .AddHttpClientInstrumentation()
                            .AddAspNetCoreInstrumentation()
                            ;

                        builder
                            .AddJaegerExporter(o => o.AgentHost = jaegerConfig.Host);
                    }
                );
        }

        return webApplicationBuilder;
    }
}