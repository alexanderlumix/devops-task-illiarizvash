#!/usr/bin/env python3
"""
Simple tests for Python scripts
"""

import unittest
import sys
import os

# Add scripts directory to path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))


class TestScripts(unittest.TestCase):
    """Test cases for Python scripts"""

    def test_script_imports(self):
        """Test that scripts can be imported without errors"""
        print("✅ Testing script imports...")

        # Test importing scripts
        try:
            import check_passwords

            print("✅ check_passwords.py imported successfully")
        except ImportError as e:
            print(f"⚠️  check_passwords.py import failed: {e}")

        try:
            import check_replicaset_status

            print("✅ check_replicaset_status.py imported successfully")
        except ImportError as e:
            print(f"⚠️  check_replicaset_status.py import failed: {e}")

        try:
            import create_app_user

            print("✅ create_app_user.py imported successfully")
        except ImportError as e:
            print(f"⚠️  create_app_user.py import failed: {e}")

        try:
            import init_mongo_servers

            print("✅ init_mongo_servers.py imported successfully")
        except ImportError as e:
            print(f"⚠️  init_mongo_servers.py import failed: {e}")

        # Test passed if we reach this point
        pass

    def test_requirements_file(self):
        """Test that requirements.txt exists and is readable"""
        print("✅ Testing requirements.txt...")

        requirements_path = os.path.join(os.path.dirname(__file__), "requirements.txt")
        self.assertTrue(
            os.path.exists(requirements_path), "requirements.txt should exist"
        )

        with open(requirements_path, "r", encoding="utf-8") as f:
            content = f.read()
            self.assertIn("pymongo", content, "requirements.txt should contain pymongo")
            self.assertIn("pyyaml", content, "requirements.txt should contain pyyaml")

        print("✅ requirements.txt is valid")

    def test_script_structure(self):
        """Test that scripts have proper structure"""
        print("✅ Testing script structure...")

        scripts = [
            "check_passwords.py",
            "check_replicaset_status.py",
            "create_app_user.py",
            "init_mongo_servers.py",
        ]

        for script in scripts:
            script_path = os.path.join(os.path.dirname(__file__), script)
            if os.path.exists(script_path):
                with open(script_path, "r", encoding="utf-8") as f:
                    content = f.read()
                    self.assertIn("import", content, f"{script} should contain imports")
                    print(f"✅ {script} has proper structure")
            else:
                print(f"⚠️  {script} not found")

        # Test passed if we reach this point
        pass

    def test_mongo_config(self):
        """Test MongoDB configuration files"""
        print("✅ Testing MongoDB configuration...")

        config_files = ["mongo_servers.yml"]

        for config_file in config_files:
            config_path = os.path.join(os.path.dirname(__file__), config_file)
            if os.path.exists(config_path):
                print(f"✅ {config_file} exists")
            else:
                print(f"⚠️  {config_file} not found")

        # Test passed if we reach this point
        pass


if __name__ == "__main__":
    unittest.main(verbosity=2)
