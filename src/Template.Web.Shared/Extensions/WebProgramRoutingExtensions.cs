using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Template.Web.Shared.Routing;

namespace Template.Web.Shared.Extensions;

public static class WebProgramRoutingExtensions {
    public static WebApplication SetRouting(this WebApplication app)
    {
        using var scope = app.Services.CreateScope();
        scope.ServiceProvider.GetServices<IApiRouter>().ToList().ForEach(apiRouter => apiRouter.AddRoutes(app));
        return app;
    }
}