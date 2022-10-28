using Template.Domain.Services;

namespace Template.Domain.Customers;

public interface ICustomers
{
    Task<Customer> GetAsync(int id);
    Task<List<Customer>> GetAsync();
    Task<Customer> Create(Customer customer);
}

public class Customers : ICustomers
{
    public async Task<Customer> GetAsync(int id) => throw new NotImplementedException();

    public async Task<List<Customer>> GetAsync() =>
        new List<Customer>
        {
            new ()
        };

    public async Task<Customer> Create(Customer customer) => throw new NotImplementedException();
}