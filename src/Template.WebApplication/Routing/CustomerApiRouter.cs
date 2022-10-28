using Template.Domain.Customers;
using Template.Domain.Services;
using Template.WebApplication.Models;

namespace Template.WebApplication.Routing;
using WebAppBuilder = Microsoft.AspNetCore.Builder;
using Functional;

public static class CustomerApiRouter
{
    private const string Route = "customer";
    private const string RouteWithId = $"{Route}/{{id}}";

    internal static WebAppBuilder.WebApplication MapCustomerEndpoints(this WebAppBuilder.WebApplication endpoints) =>
        endpoints
            .Tee(_ => _.MapGet(Route, GetAsync).DescribeGet())
            .Tee(_ => _.MapGet(RouteWithId, GetByIdAsync).DescribeGet())
            .Tee(_ => _.MapPost(Route, PostAsync).DescribePost());

    private static readonly Delegate GetAsync =
        (ICustomers customers) => customers.GetAsync();

    private static readonly Delegate PostAsync = 
        async (ICustomers customers, CustomerModel model, CustomerValidator validator) =>
        {
            // TODO: Consider use of Decorator pattern to make an onion 
            // consisting on three steps:
            // 1. Validate
            // 2. Map
            // 3. Action invocation
            
            var result = validator.Validate(model);
            if (!result.IsValid)
                return Results.ValidationProblem(result.ToDictionary());
            
            var customer = await model
                .Map(CustomerMapper.ToEntity)
                .Tee(customers.Create);
            
            return Results.Created($"/{customer.Id}", customer);
        };
    
    private static readonly Delegate GetByIdAsync =
        (ICustomers customers, int id) => customers.GetAsync(id);

    private static RouteHandlerBuilder DescribeGet(this RouteHandlerBuilder route) =>
        route.Produces(StatusCodes.Status200OK, typeof(Customer))
             .Produces(StatusCodes.Status400BadRequest, typeof(ErrorResponse));
    
    private static RouteHandlerBuilder DescribePost(this RouteHandlerBuilder route) =>
        route.Produces(StatusCodes.Status200OK, typeof(Customer))
             .Produces(StatusCodes.Status400BadRequest, typeof(ErrorResponse));
}