### Table of Contents

- [GORM Relation Has Many](#gorm-relation-has-many)
  - [Repository](#repository)

---

# GORM Relation Has Many

Reference: [Official GORM Website](https://gorm.io/docs/has_many.html)


## Models

- Don't forget to add  relation has many in `user.go` models like this :
  ```go
  type User struct {
    ID        int                   `json:"id"`
    Name      string                `json:"name" gorm:"type: varchar(255)"`
    Email     string                `json:"email" gorm:"type: varchar(255)"`
    Password  string                `json:"-" gorm:"type: varchar(255)"`
    Profile   ProfileResponse       `json:"profile"`
    Products  []ProductUserResponse `json:"products"`
    CreatedAt time.Time             `json:"-"`
    UpdatedAt time.Time             `json:"-"`
  }
  ```

## Relation

For this section, example Has Many relation:

- `User` &rarr; `Product`: to get User Product

## Repository

- Inside `repositories` folder, in `user.go` file write this below code

  > File: `repositories/user.go`

  ```go
  func (r *repository) FindUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.Preload("Profile").Preload("Products").Find(&users).Error // add this code

    return users, err
  }

  func (r *repository) GetUser(ID int) (models.User, error) {
    var user models.User
    err := r.db.Preload("Profile").Preload("Products").First(&user, ID).Error // add this code

    return user, err
  }
  ```

  \*In this case, just add `Preload` to make relation
