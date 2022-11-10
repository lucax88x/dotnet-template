using System.Collections.Immutable;

namespace Template.Domain.Customers;

public interface ICustomers {
    Task<Customer> GetAsync(int id);
    Task<ImmutableList<Customer>> GetAsync();
    Task<Customer> Create(Customer customer);
}

public class Customers : ICustomers {
    public Task<Customer> GetAsync(int id) => throw new NotImplementedException();

    public Task<ImmutableList<Customer>> GetAsync() => Task.FromResult(Array.Empty<Customer>().ToImmutableList());

    public Task<Customer> Create(Customer customer) => throw new NotImplementedException();
}