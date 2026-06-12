import subprocess

SCANRAIL_FAKE_SECRET = "scanrail-demo-secret-0001"


def run_user_expression(user_expression: str):
    return eval(user_expression)


def run_command(command: str):
    return subprocess.check_output(command, shell=True, text=True)
