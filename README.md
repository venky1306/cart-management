## Cart Management Web App Backend

This backend application is built using Go with Gin framework, MongoDB for data storage, JWT for authentication and authorization features.

### Features

- **User Authentication**: Endpoints for user signup and login with JWT token generation.
- **Product Management**: APIs to search for products, add new products, and manage cart items.
- **Cart Operations**: Functions to add, remove, checkout, and instantly buy items from the cart.

### Endpoints

#### User Routes
- **POST /v1/users/signup**: Register a new user.
- **POST /v1/users/login**: User login with credentials.

#### Product Routes
- **GET /v1/products**: Retrieve all products.
- **GET /v1/products/search**: Search for products by query.
- **POST /v1/products/addproduct**: Add a new product to the inventory.

#### Cart Operations
- **GET /v1/addtocart**: Add an item to the cart.
- **GET /v1/removeitem**: Remove an item from the cart.
- **GET /v1/checkout**: Proceed to checkout items in the cart.
- **GET /v1/instantbuy**: Instantly buy an item from the cart.

### Running Locally

To run this application locally, follow these steps:

1. **Start a MongoDB Container**
    ```shell
    # Run MongoDB container
    docker run --name my-mongodb-container -d -p 27017:27017 mongo
    ```

2. **Build and Run the Application**
    ```shell
    # Build the Docker image
    docker build -t cart-management .

    # Run the Docker container
    docker run --name goserver --publish 8001:8000 --env MONGODB_URL="mongodb://<YOUR_MONGODB_IP>:27017/" --env SECRET_KEY="<YOUR_SECRET_KEY>"cart-management
    ```
    Replace `<YOUR_MONGODB_IP>` with your MongoDB server's container IP address.
    Replace `<YOUR_SECRET_KEY>` with your JWT SECRET.

    Make sure to run both mongodb container and goserver in same docker network.

### External Go Packages Used

- **github.com/gin-gonic/gin**: Gin framework for building web applications in Go.
- **go.mongodb.org/mongo-driver**: MongoDB driver for Go, enabling interaction with MongoDB.
- **github.com/golang-jwt/jwt/v5**: JWT implementation for Go, facilitating authentication with JSON Web Tokens.
