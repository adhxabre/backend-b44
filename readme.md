# Insert data with ORM

### Repositories

> File: `repositories/user.go`

- Add User data using `Create` method

  ```go
  func (r *repository) CreateUser(user models.User) (models.User, error) {
    err := r.db.Create(&user).Error // Using Create method

    return user, err
  }
  ```
