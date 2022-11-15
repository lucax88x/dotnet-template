using Microsoft.AspNetCore.Builder;

namespace Template.Web.Common.Extensions;

public static class WebProgramAuthExtensions {
    public static WebApplicationBuilder AddTemplateAuth(this WebApplicationBuilder builder)
    {
        // builder.Services
        //     .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
        //     .AddMicrosoftIdentityWebApi(builder.Configuration.GetSection("AzureAd"));

        return builder;
    }

    public static WebApplication UseTemplateAuth(this WebApplication app)
    {
        app.UseAuthentication();
        app.UseAuthorization();
        return app;
    }
}