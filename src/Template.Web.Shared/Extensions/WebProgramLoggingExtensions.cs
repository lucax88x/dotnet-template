using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Logging;

namespace Template.Web.Shared.Extensions;

public static class WebProgramLoggingExtensions {
    public static WebApplicationBuilder SetLogging(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Logging.ClearProviders();
        webApplicationBuilder.Logging.AddConsole();
        return webApplicationBuilder;
    }
}