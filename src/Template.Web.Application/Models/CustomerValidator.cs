﻿using FluentValidation;

namespace Template.Web.Application.Models;

public class CustomerValidator : AbstractValidator<CustomerModel>
{
    public CustomerValidator()
    {
        RuleFor(x => x.FirstName).NotEmpty().Length(3, 50);
        RuleFor(x => x.LastName).NotEmpty().Length(3, 50);
    }
}