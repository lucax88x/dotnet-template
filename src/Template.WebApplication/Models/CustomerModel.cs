using Template.Domain.Services;

namespace Template.WebApplication.Models;

public record CustomerModel(
    int CustomerId,
    string FirstName,
    string LastName,
    string CompanyName,
    string EmailAddress,
    ICollection<string> PhoneNumbers
)
{
    public Customer ToEntity()
    {
        throw new NotImplementedException();
    }
}