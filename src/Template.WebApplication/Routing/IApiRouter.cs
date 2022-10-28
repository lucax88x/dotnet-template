namespace Template.WebApplication.Routing;

public interface IApiRouter
{
    public Microsoft.AspNetCore.Builder.WebApplication AddRoutes(Microsoft.AspNetCore.Builder.WebApplication app);
}