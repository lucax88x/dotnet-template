using Template.Domain.Customers;

namespace Template.Web.Application;

public static class ProgramExtensions {
    public static WebApplicationBuilder RegisterActions(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Services.AddTransient<ICustomerService, CustomerServiceService>();
        return webApplicationBuilder;
    }
}