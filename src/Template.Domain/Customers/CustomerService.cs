using System.Collections.Immutable;
using System.Diagnostics;
using Template.Common;
using Template.Common.Tracing;

namespace Template.Domain.Customers;

public interface ICustomerService {
    Task<Customer> GetAsync(int id);
    Task<ImmutableList<Customer>> GetAsync();
    Task<Customer> Create(Customer customer);
}

public class CustomerService : ICustomerService {
    public Task<Customer> GetAsync(int id) => throw new NotImplementedException();

    public Task<ImmutableList<Customer>> GetAsync()
    {
        using var activity = Tracing.Source.StartActivity();

        return Task.FromResult(
            new List<Customer> { new() { Id = 1, Name = "some name" }, new() { Id = 2, Name = "another name" } }.ToImmutableList()
        );
    }

    public Task<Customer> Create(Customer customer) => throw new NotImplementedException();
}