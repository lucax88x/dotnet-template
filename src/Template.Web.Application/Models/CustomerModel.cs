using Template.Domain.Customers;

namespace Template.Web.Application.Models;

public record CustomerModel(
    int CustomerId,
    string FirstName,
    string LastName,
    string CompanyName,
    string EmailAddress,
    IReadOnlyCollection<string> PhoneNumbers
)
{
    public Customer ToEntity()
    {
        throw new NotImplementedException();
    }
}