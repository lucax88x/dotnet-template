using Functional;
using Template.Domain.Customers;
using Template.Web.Application.Models;

namespace Template.Web.Application.Routing;

using WebAppBuilder = Microsoft.AspNetCore.Builder;

public static class CustomerApiRouter {
    const string _route = "customer";
    const string _routeWithId = $"{_route}/{{id}}";

    internal static WebApplication MapCustomerEndpoints(this WebApplication endpoints) =>
        endpoints
            .Tee(_ => _.MapGet(_route, GetAsync).DescribeGet())
            .Tee(_ => _.MapGet(_routeWithId, GetByIdAsync).DescribeGet())
            .Tee(_ => _.MapPost(_route, PostAsync).DescribePost());

    static readonly Delegate GetAsync =
        (ICustomerService customers) => customers.GetAsync();

    static readonly Delegate PostAsync =
        async (ICustomerService customers, CustomerModel model, CustomerValidator validator) =>
        {
            // TODO: Consider use of Decorator pattern to make an onion 
            // consisting on three steps:
            // 1. Validate
            // 2. Map
            // 3. Action invocation

            var result = validator.Validate(model);
            if (!result.IsValid) return Results.ValidationProblem(result.ToDictionary());

            var customer = await model
                .Map(CustomerMapper.ToEntity)
                .Tee(customers.Create);

            return Results.Created(new Uri($"/{customer.Id}"), customer);
        };

    static readonly Delegate GetByIdAsync =
        (ICustomerService customers, int id) => customers.GetAsync(id);

    static RouteHandlerBuilder DescribeGet(this RouteHandlerBuilder route) =>
        route.Produces(StatusCodes.Status200OK, typeof(Customer))
            .Produces(StatusCodes.Status400BadRequest, typeof(ErrorResponse));

    static RouteHandlerBuilder DescribePost(this RouteHandlerBuilder route) =>
        route.Produces(StatusCodes.Status200OK, typeof(Customer))
            .Produces(StatusCodes.Status400BadRequest, typeof(ErrorResponse));
}