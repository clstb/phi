"""The command-line interface for Fava."""
from __future__ import annotations
import os
import click
from werkzeug.middleware.profiler import ProfilerMiddleware
from fava.application import app


@click.option("-env", "--environment", help="In which env to run flask app.")
@click.option("-d", "--debug", is_flag=True, help="Turn on debugging.")
@click.option("--profile", is_flag=True, help="Turn on profiling. " "Implies --debug.")
@click.option(
    "--profile-dir",
    type=click.Path(),
    help="Output directory for profiling data.",
)
@click.version_option(prog_name="fava")
def main(env="development", debug=True, profile=False, profile_dir=None):
    os.environ["FLASK_ENV"] = env
    if profile:
        app.config["PROFILE"] = True
    if profile:
        app.wsgi_app = ProfilerMiddleware(  # type: ignore
            app.wsgi_app,
            restrictions=(30,),
            profile_dir=profile_dir if profile_dir else None,
        )

    app.jinja_env.auto_reload = os.environ["FLASK_ENV"] == "development"
    app.run("0.0.0.0", 8083, debug=debug)


if __name__ == "__main__":
    app.config["BEANCOUNT_FILES"] = []
    main()
