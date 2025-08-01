# MongoDB Application User Creation Script
# This script creates a dedicated user for the application with read/write permissions
import pymongo

# MongoDB connection parameters
MONGO_HOST = '127.0.0.1'
MONGO_PORT = 27031
ADMIN_USER = 'mongo-1'
ADMIN_PASS = 'mongo-1'

# Application database and user configuration
APP_DB = 'appdb'
APP_USER = 'appuser'
APP_PASS = 'appuserpassword'


def create_app_user():
    """Create application user with read/write permissions on appdb database"""
    # Connect to MongoDB using admin credentials
    uri = f"mongodb://{ADMIN_USER}:{ADMIN_PASS}@{MONGO_HOST}:{MONGO_PORT}/admin?replicaSet=rs0&authSource=admin"
    client = pymongo.MongoClient(uri, serverSelectionTimeoutMS=5000)
    db = client[APP_DB]
    
    try:
        # Create user with readWrite role on the application database
        db.command("createUser", APP_USER,
                   pwd=APP_PASS,
                   roles=[{"role": "readWrite", "db": APP_DB}])
        print(f"User '{APP_USER}' created with readWrite role on database '{APP_DB}'.")
    except pymongo.errors.OperationFailure as e:
        if 'already exists' in str(e):
            print(f"User '{APP_USER}' already exists.")
        else:
            print(f"Failed to create user: {e}")
    finally:
        client.close()

if __name__ == '__main__':
    create_app_user()
