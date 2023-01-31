using Template.Domain.Customers;

namespace Template.Domain.Tests;

public class UnitTest1
{
    [Fact]
    public void ShouldBeTrue()
    {
        var customer = new Customer() { Name = "name", Id = 1 };
        
        Assert.Equal("name", customer.Name);
    }
    
    [Fact]
    public void ShouldBeFalse()
    {
        var customer = new Customer() { Name = "name", Id = 1 };
        
        Assert.Equal("name2", customer.Name);
    }
}