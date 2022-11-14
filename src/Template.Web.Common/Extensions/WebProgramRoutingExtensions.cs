using System.Linq;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Template.Web.Common.Routing;

namespace Template.Web.Common.Extensions;

public static class WebProgramRoutingExtensions {
    public static WebApplication UseTemplateRouting(this WebApplication app)
    {
        using var scope = app.Services.CreateScope();
        scope.ServiceProvider.GetServices<IApiRouter>().ToList().ForEach(apiRouter => apiRouter.AddRoutes(app));
        return app;
    }
}