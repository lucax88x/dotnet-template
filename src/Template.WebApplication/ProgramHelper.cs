using Template.Domain.Customers;
using Template.WebApplication.Routing;

namespace Template.WebApplication;

public static class ProgramHelper
{
    public static Microsoft.AspNetCore.Builder.WebApplication SetAuth(this Microsoft.AspNetCore.Builder.WebApplication webApplication)
    {
        webApplication.UseAuthentication();
        webApplication.UseAuthorization();
        return webApplication;
    }

    public static Microsoft.AspNetCore.Builder.WebApplication SetSwaggerWhenDevelopment(this Microsoft.AspNetCore.Builder.WebApplication app)
    {
        if (!app.Environment.IsDevelopment()) return app;
        app.UseSwagger();
        app.UseSwaggerUI();
        return app;
    }

    public static Microsoft.AspNetCore.Builder.WebApplication SetRouting(this Microsoft.AspNetCore.Builder.WebApplication app)
    {
        using var scope = app.Services.CreateScope();
        scope.ServiceProvider.GetServices<IApiRouter>().ToList().ForEach(
            apiRouter => apiRouter.AddRoutes(app));
        return app;
    }

    public static WebApplicationBuilder SetLogging(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Logging.ClearProviders();
        webApplicationBuilder.Logging.AddConsole();
        return webApplicationBuilder;
    }
    
    public static WebApplicationBuilder RegisterActions(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Services.AddTransient<ICustomers, Customers>();
        return webApplicationBuilder;
    }
    
    public static WebApplicationBuilder SetApiExplorerAndSwagger(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Services.AddEndpointsApiExplorer()
            .AddSwaggerGen();
        return webApplicationBuilder;
    }
}