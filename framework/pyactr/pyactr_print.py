"""
pyactr_print adds some extra print capabilities to pyactr productions.

pyactr is limited in what it can print using "show" to one named slot of the buffer.

    !goal>
        show start

pyactr_print adds the ability to print strings, numbers, buffer contents, and slots (by name)
from multiple buffers. It uses a new buffer called "print" with commands "text" and "buffer":

	!print>
	    text "'Start is ', goal.start, ' and second is ', 'retrieval.second'"
        buffer goal
        buffer retrieval.word

This works by using fields from ACTRModel (reads & modifies "__buffers") and
Buffer (reads "_data") directly since there are no accessors for them.

Unfortunately due to the way pyactr is implemented, we are currently limited to
one print command per production. It does, however, allow multiple "text"s in
one command:

	!print>
	    text "'a string'"
	    text "42"

To use this new buffer, construct it passing in the model like this:

    import pyactr_print

    pyactr_print.PrintBuffer(pyactr_fan)

Then you can use it with "!print>" as shown above.
"""

# We use csv to parse the print text we are generating.
# This is just simpler than writing it ourselves(i.e. handling "foo, bar ", 66).
import csv
import pyactr as actr

from pyactr.buffers import Buffer


class PrintBuffer(actr.buffers.Buffer):
    def __init__(self, model: actr.ACTRModel):
        actr.buffers.Buffer.__init__(self, None, None)
        self.ACTR_MODEL = model
        model._ACTRModel__buffers["print"] = self

    def text(self, *args):
        """
        Provides the "text" command.
        Prints the args - including strings, numbers, buffer contents, and slots (by name).

        Examples:
            !print>
                text "'Start is ', goal.start, ' and second is ', 'retrieval.second'"
                text "'retrieval contents: ', retrieval"
                text "'a string'"
                text "42"
        """
        text = "".join(args[1:]).strip('"')
        output = ""  # build up our output in this buffer

        for itemlist in csv.reader([text]):
            for item in itemlist:
                item = item.strip(" ")

                # Handle string
                if item[0] == "'" or item[0] == '"':
                    output += item[1:-1]
                else:
                    # Handle number
                    try:
                        float(item)
                        output += item
                    except ValueError:
                        # If we are here, we should have a buffer or a buffer.slotname
                        output += self.get_buffer_data(item)

        print(output)

    def buffer(self, *args):
        """
        Provides the "buffer" command.
        Prints the contents of a buffer or a buffer's slot.

        Examples:
            !print>
                buffer retrieval
                buffer retrieval.word
        """
        name = "".join(args)
        contents = self.get_buffer_data(name)
        print(f"{name}: {contents}")

    def get_buffer_data(self, item: str) -> str:
        """
        Given an "item" which is either a <buffer name> or a <buffer name>.<slot name>,
        return the contents.
        """
        ids = item.split(".")
        match len(ids):
            case 1:
                contents = self.get_buffer_contents(item)
            case 2:
                contents = self.get_slot_contents(ids[0], ids[1])
            case _:
                print(
                    "ERROR: expected <buffer> or <buffer>.<slot_name>, found '"
                    + item
                    + "'"
                )
                raise KeyError
        return contents

    def get_buffer_contents(self, buffer_name: str) -> str:
        """
        Gets all the contents of a buffer.
        """
        buffer = self.get_buffer(buffer_name)
        data = buffer._data

        if data:
            return str(data.copy().pop())
        else:
            return "<empty>"

    def get_slot_contents(self, buffer_name: str, slot_name: str) -> str:
        """
        Gets the contents of a specific buffer slot.
        """
        buffer = self.get_buffer(buffer_name)

        if buffer._data:
            chunk = buffer._data.copy().pop()
        else:
            chunk = None

        try:
            return str(getattr(chunk, slot_name))
        except AttributeError:
            print(
                "ERROR: no slot named '"
                + slot_name
                + "' in buffer '"
                + buffer_name
                + "'"
            )
            raise

    def get_buffer(self, buffer_name: str) -> Buffer:
        """
        Gets a buffer by name and returns it.
        """
        if buffer_name in self.ACTR_MODEL._ACTRModel__buffers:
            return self.ACTR_MODEL._ACTRModel__buffers[buffer_name]

        print("ERROR: Buffer '" + buffer_name + "' not found")
        raise KeyError

    def add(self, *args):
        raise AttributeError(
            "Attempt to add an element to the print buffer. This is not allowed."
        )

    def clear(self, *args):
        raise AttributeError("Attempt to clear the print buffer. This is not allowed.")

    def create(self, *args):
        raise AttributeError(
            "Attempt to add an element to the print buffer. This is not allowed."
        )

    def retrieve(self, *args):
        raise AttributeError(
            "Attempt to retrieve from the print buffer. This is not allowed."
        )

    def test(self, *args):
        raise AttributeError(
            "Attempt to test the print buffer state. This is not allowed."
        )
