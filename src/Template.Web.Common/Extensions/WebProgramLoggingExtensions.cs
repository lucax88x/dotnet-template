using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Serilog;
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
        Log.Logger = LoggerConfigurationBuilder
            .BuildForHost()
            .CreateBootstrapLogger();

        hostBuilder.UseSerilog(
            (_, services, loggerConfiguration) =>
            {
                var seqConfig = services.GetRequiredService<IOptions<SeqConfig>>();

                LoggerConfigurationBuilder.BuildForApplication(new LogOptions { SeqConfig = seqConfig.Value }, loggerConfiguration);
            }
        );

        return hostBuilder;
    }
}