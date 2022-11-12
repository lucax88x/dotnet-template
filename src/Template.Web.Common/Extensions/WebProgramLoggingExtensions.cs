using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Serilog;

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
            (context, services, loggerConfiguration) =>
            {
                LoggerConfigurationBuilder.BuildForApplication(loggerConfiguration);
            }
        );

        return hostBuilder;
    }
}