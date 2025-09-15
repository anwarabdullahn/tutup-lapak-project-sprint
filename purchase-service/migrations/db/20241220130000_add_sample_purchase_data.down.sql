-- Remove sample data
DELETE FROM purchase_senders WHERE id IN (
    '01234567-89ab-cdef-0123-456789abcde5',
    '01234567-89ab-cdef-0123-456789abcde6'
);

DELETE FROM purchase_items WHERE id IN (
    '01234567-89ab-cdef-0123-456789abcde2',
    '01234567-89ab-cdef-0123-456789abcde3',
    '01234567-89ab-cdef-0123-456789abcde4'
);

DELETE FROM purchases WHERE id IN (
    '01234567-89ab-cdef-0123-456789abcdef',
    '01234567-89ab-cdef-0123-456789abcde1'
);
