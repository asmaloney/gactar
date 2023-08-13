"""
pyactr_print adds some extra print capabilities to pyactr productions.

pyactr is limited in what it can print using "show" to one named slot of the buffer.

    !goal>
        show start

pyactr_print adds the ability to print strings, numbers, and slots (by name) from
multiple buffers by patching ACTRModel and Buffer. It uses a new command "print_text"
in productions which takes a string:

	!goal>
	    print_text "'Start is ', goal.start, ' and second is ', 'retrieval.second'"

Unfortunately due to the way pyactr is implemented, we are currently limited to
one "print"text" statement per production.
"""

# We use csv to parse the print text we are generating.
# This is just simpler than writing it ourselves(i.e. handling "foo, bar ", 66).
import csv
import pyactr as actr

from pyactr.buffers import Buffer

# This is generally Bad(tm). This means we can only have one actr.ACTRModel.
# Unfortunately, there's no way to get from a Buffer to an ACTRModel.
ACTR_MODEL: actr.ACTRModel


def set_model(model: actr.ACTRModel):
    """
    Sets our module's ACTR_MODEL so we can access it in Buffer.print_text.
    """
    global ACTR_MODEL
    ACTR_MODEL = model


def get_buffer(self, buffer_name: str) -> Buffer:
    """
    Gets a buffer from a name and returns it.
    """
    if buffer_name in self._ACTRModel__buffers:
        return self._ACTRModel__buffers[buffer_name]

    print('ERROR: Buffer \'' + buffer_name + '\' not found')
    raise KeyError


# Monkey patch ACTRModel to add a new method.
actr.ACTRModel.get_buffer = get_buffer


def get_slot_contents(self, buffer_name: str, slot_name: str) -> str:
    """
    Gets the contents of a slot.
    """
    if self._data:
        chunk = self._data.copy().pop()
    else:
        chunk = None

    try:
        return str(getattr(chunk, slot_name))
    except AttributeError:
        print('ERROR: no slot named \'' + slot_name +
              '\' in buffer \'' + buffer_name + '\'')
        raise


def print_text(*args):
    """
    Prints the args - including strings, numbers, and slots (by name).
    """
    text = ''.join(args[1:]).strip('"')
    output = ''  # build up our output in this buffer

    for itemlist in csv.reader([text]):
        for item in itemlist:
            item = item.strip(' ')

            # Handle string
            if item[0] == '\'' or item[0] == '"':
                output += item[1:-1]
            else:
                # Handle number
                try:
                    float(item)
                    output += item
                except ValueError:
                    # If we are here, we should have a buffer.slotname
                    ids = item.split('.')
                    if len(ids) != 2:
                        print(
                            'ERROR: expected <buffer>.<slot_name>, found \'' +
                            item + '\'')
                    else:
                        buffer = ACTR_MODEL.get_buffer(ids[0])
                        output += buffer.get_slot_contents(ids[0], ids[1])

    print(output)


# Monkey patch Buffer to add a new methods.
Buffer.get_slot_contents = get_slot_contents
Buffer.print_text = print_text
