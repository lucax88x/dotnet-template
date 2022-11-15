using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace Template.Web.Common.Extensions;

public static class WebProgramSwaggerExtensions {
    public static WebApplicationBuilder AddSwaggerWhenDevelopment(this WebApplicationBuilder builder)
    {
        if (!builder.Environment.IsDevelopment()) return builder;

        builder
            .Services
            .AddEndpointsApiExplorer()
            .AddSwaggerGen();

        return builder;
    }

    public static WebApplication UseSwaggerWhenDevelopment(this WebApplication app)
    {
        if (!app.Environment.IsDevelopment()) return app;

        app.UseSwagger();
        app.UseSwaggerUI();

        app.Logger.LogInformation("Swagger available at /swagger");

        return app;
    }
}