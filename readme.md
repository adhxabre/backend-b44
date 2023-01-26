### Table of Contents

- [GORM Relation belongs to](#gorm-relation-has-one)
  - [Repository](#repository)

---

# GORM Relation Has One

Reference: [Official GORM Website](https://gorm.io/docs/has_one.html)

## Relation

For this section, example Has One relation:

- `User` &rarr; `Profile`: to get User Profile

## Models

- Don't forget to add `Profile`, so User can access Profile response :

  ```go
  type User struct {
    ID        int             `json:"id"`
    Name      string          `json:"name" gorm:"type: varchar(255)"`
    Email     string          `json:"email" gorm:"type: varchar(255)"`
    Password  string          `json:"-" gorm:"type: varchar(255)"`
    Profile   ProfileResponse `json:"profile"`
    CreatedAt time.Time       `json:"-"`
    UpdatedAt time.Time       `json:"-"`
  }
  ```

## Repository

- Inside `repositories` folder, in `user.go` file write this below code

  > File: `repositories/user.go`

  ```go
  func (r *repository) FindUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.Preload("Profile").Find(&users).Error // add this code

    return users, err
  }

  func (r *repository) GetUser(ID int) (models.User, error) {
    var user models.User
    err := r.db.Preload("Profile").First(&user, ID).Error // add this code

    return user, err
  }
  ```

  \*In this case, just add `Preload` to make relation
