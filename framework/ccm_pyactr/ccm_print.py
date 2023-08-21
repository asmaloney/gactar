"""
ccm_print adds some extra print capabilities to ccm productions.
"""

from typing import Dict, List

from python_actr import Buffer


class CCMPrint():
    def __init__(self):
        self.chunk_map: Dict[str, List[str]] = {}

    def register_chunk(self, chunk_name: str, slot_names: List[str]):
        """
        Registers a chunk name with the slots which are available to it.
        """
        self.chunk_map[chunk_name] = slot_names

    def print_chunk(self,  buffer: Buffer, buffer_name: str, ):
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
        item = buffer.__getitem__(slot_index + 1)

        print(f"{buffer_name}.{slot}: {item}")
