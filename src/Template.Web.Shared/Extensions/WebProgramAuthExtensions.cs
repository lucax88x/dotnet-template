using Microsoft.AspNetCore.Builder;

namespace Template.Web.Shared.Extensions;

public static class WebProgramAuthExtensions {
    public static WebApplication SetAuth(this WebApplication webApplication)
    {
        webApplication.UseAuthentication();
        webApplication.UseAuthorization();
        return webApplication;
    }
}