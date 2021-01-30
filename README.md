# xmlcheck-ppss-trans
 Check/repair the XML of a PPSS Transaction input file

Use case is that I had a bunch of existing XML files
of "unknown" correctness (as in: some were only half complete)
that also needed a new element/value added that fortunately
had to match an existing value.

So, "invalid" XML files are suffixed as .BAD to keep them
from being picked up and to make them easy to spot for a
human who is looking for them.

Then the value of Day is copied into WK_DAY to save a bunch
of manual ETL re-work.

Finally, the big, ugly struct for Record is making most of you
SMH right now. I didn't need the full struct, but I wanted to
get the practice of being complete about it. And, yes, this is
how it looks in the XSD (caps and underscores which the linter
absolutely LOVES). So IIWII (It Is What It Is).

And it bears saying: the app loading this cannot tolerate
ANY variance from its XSD. It ain't complicated, but it also
ain't flexible.
