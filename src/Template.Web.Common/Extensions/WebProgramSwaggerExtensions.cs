using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

namespace Template.Web.Common.Extensions;

public static class WebProgramSwaggerExtensions {
    public static WebApplicationBuilder AddSwaggerWhenDevelopment(this WebApplicationBuilder builder)
    {
        if (!builder.Environment.IsDevelopment()) return builder;

        builder
            .Services
            .AddSwaggerGen();

        return builder;
    }

    public static WebApplicationBuilder AddApiExplorer(this WebApplicationBuilder builder)
    {
        builder
            .Services
            .AddEndpointsApiExplorer();

        return builder;
    }

    public static WebApplication UseSwaggerWhenDevelopment(this WebApplication app)
    {
        if (!app.Environment.IsDevelopment()) return app;

        app.UseSwagger();
        app.UseSwaggerUI();

        return app;
    }
}