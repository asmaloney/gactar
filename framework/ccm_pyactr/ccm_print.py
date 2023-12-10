"""
ccm_print adds some extra print capabilities to ccm productions.

To use it, create an instance and regsiter each chunk with it:

    printer = CCMPrint()
    printer.register_chunk("add", ["num1", "num2", "count", "sum"])
    printer.register_chunk("count", ["number", "next"])

Then in the productions:
    printer.print_chunk_slot(goal, "goal", "sum")
"""

from typing import Dict, List

from python_actr import Buffer


class CCMPrint:
    """
    CCMPrint provides methods to print the contents of a buffer or the contents
    of one slot of a buffer.
    """

    def __init__(self):
        self.chunk_map: Dict[str, List[str]] = {}

    def register_chunk(self, chunk_name: str, slot_names: List[str]):
        """
        Registers a chunk name with the slots which are available to it.
        """
        self.chunk_map[chunk_name] = slot_names

    def print_chunk(
        self,
        buffer: Buffer,
        buffer_name: str,
    ):
        """
        Prints the contents of the buffer.

        e.g. retrieval: word Mary ProperN
        """
        print(f"{buffer_name}: {buffer.chunk}")

    def print_chunk_slot(self, buffer: Buffer, buffer_name: str, slot: str):
        """
        Prints the contents of one slot of the buffer.

        e.g. retrieval.form: Mary
        """
        # get the slot list from the chunk type
        slot_list = self.chunk_map[buffer.chunk[0]]

        # find the slot index so we can look up the contents
        slot_index = slot_list.index(slot)

        # +1 because the first item is the chunk type
        item = buffer[slot_index + 1]

        print(f"{buffer_name}.{slot}: {item}")
