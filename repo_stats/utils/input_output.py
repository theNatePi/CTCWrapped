from termcolor import colored

class InputOutputHandler:
    """Handle input and output."""
    def __init__(self, verbose: bool = False):
        if verbose:
            self.outputter = Outputter()
        else:
            self.outputter = DummyOutputter()
        
    def output(self, title: str, message: str, message_type: str) -> None:
        self.outputter.output(title, message, message_type)


class DummyOutputter:
    """Dummy outputter that does nothing."""
    def output(self, title: str, message: str, message_type: str) -> None:
        pass

    def error(self, title: str, message: str) -> None:
        pass


class Outputter(DummyOutputter):
    """Output messages."""
    def output(self, title: str, message: str, message_type: str) -> None:
        """Output a formatted message with color based on message type.
        Message type should be out of ['success', 'progress', 'error']"""
        color = 'white'
        if message_type == "success":
            color = 'green'
        elif message_type == "progress":
            color = 'light_blue'
        elif message_type == "error":
            color = 'red'

        open_bracket = colored('[', attrs = ['dark'])
        close_bracket = colored(']', attrs = ['dark'])
        print(f"{open_bracket}"
              f"{colored(self._current_request_type, 'light_green')}{close_bracket} "
              f"{open_bracket}{colored(self.rate_limit_remaining, 'light_magenta', attrs = ['dark'])}"
              f"{close_bracket} {colored(title, color)}: "
              f"{colored(message, attrs = ['dark'])}")
    
    def error(self, title: str, message: str) -> None:
        self.output(title, message, "error")


class PrintWrapper:
    """Print wrapper."""
    def output(self, title: str, message: str, message_type: str) -> None:
        print(message_type, title, message)
