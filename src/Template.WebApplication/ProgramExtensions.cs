using Template.Domain.Customers;

namespace Template.WebApplication;

public static class ProgramExtensions {
    public static WebApplicationBuilder RegisterActions(this WebApplicationBuilder webApplicationBuilder)
    {
        webApplicationBuilder.Services.AddTransient<ICustomers, Customers>();
        return webApplicationBuilder;
    }
}