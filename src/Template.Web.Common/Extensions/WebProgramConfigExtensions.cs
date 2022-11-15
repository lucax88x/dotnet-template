using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

namespace Template.Web.Common.Extensions;

public static class WebProgramConfigExtensions {
    public static WebApplicationBuilder AddAndValidateOptions<TConfig>(this WebApplicationBuilder builder, string section)
        where TConfig : class
    {
        builder.Services
            .AddOptions<TConfig>()
            .Bind(builder.Configuration.GetSection(section))
            .ValidateDataAnnotations()
            .ValidateOnStart();

        return builder;
    }
}