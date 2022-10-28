using Template.Domain.Services;

namespace Template.WebApplication.Models;

public static class CustomerMapper
{
    public static Customer ToEntity(CustomerModel model) =>
        new ()
        {
            Id = DateTime.UtcNow.Ticks,
            Name = $"{model.LastName}, {model.FirstName}"
        };
}