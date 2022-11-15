using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Serilog;
using Serilog.Sinks.FastConsole;
using Template.Common;
using Template.Web.Common.Configs;

namespace Template.Web.Common.Extensions;

public static class WebProgramLoggingExtensions {
    public static WebApplicationBuilder AddTemplateLogging(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Logging.ClearProviders();
        webApplicationBuilder.Host.UseTemplateLogging();

        return webApplicationBuilder;
    }

    public static IHostBuilder UseTemplateLogging(this IHostBuilder hostBuilder)
    {
        Log.Logger = new LoggerConfiguration()
            .WriteTo.FastConsole()
            .CreateBootstrapLogger();

        hostBuilder.UseSerilog(
            (_, services, loggerConfiguration) =>
            {
                var seqConfig = services.GetService<IOptions<SeqConfig>>();

                loggerConfiguration
                    .Enrich.FromLogContext()
                    .Enrich.WithProperty("source", SourceAndVersion.SourceName)
                    .WriteTo.FastConsole();

                if (seqConfig is not null)
                {
                    loggerConfiguration.WriteTo.Seq(seqConfig.Value.Host);
                }
            }
        );

        return hostBuilder;
    }
}