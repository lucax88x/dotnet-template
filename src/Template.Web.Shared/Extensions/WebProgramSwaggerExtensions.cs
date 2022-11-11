using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace Template.Web.Shared.Extensions;

public static class WebProgramSwaggerExtensions {
    public static WebApplicationBuilder SetApiExplorerAndSwagger(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Services.AddEndpointsApiExplorer()
            .AddSwaggerGen();
        return webApplicationBuilder;
    }

    public static WebApplication SetSwaggerWhenDevelopment(this WebApplication app)
    {
        if (!app.Environment.IsDevelopment()) return app;
        app.UseSwagger();
        app.UseSwaggerUI();
        return app;
    }
}