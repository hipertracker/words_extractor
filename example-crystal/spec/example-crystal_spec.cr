require "./spec_helper"

describe Example::Crystal do
  it "should be understand PL collations" do
    words = %w(ala ąla ćma ciemno ęle element łódź lody śmierć serial źdźbło żółw zawias)
    expected = %w(ala ąla ciemno ćma element ęle lody łódź serial śmierć zawias źdźbło żółw)
    words.sort { |x, y| Example::Crystal.word_cmp(x, y) }.should eq(expected)
  end

  it "should be understand PL collations for q" do
    words = %w(ala ąla ćma ciemno ęle element łódź lody śmierć serial źdźbło żółw zawias querty)
    expected = %w(ala ąla ciemno ćma element ęle lody łódź querty serial śmierć zawias źdźbło żółw)
    words.sort { |x, y| Example::Crystal.word_cmp(x, y) }.should eq(expected)
  end

  it "should understand PL collations for upper chars" do
    words = %w(ala ąla Ćma ciemno Ęle Element łódź lody śmierć serial źdźbło żółw zawias querty)
    expected = %w(ala ąla ciemno Ćma Element Ęle lody łódź querty serial śmierć zawias źdźbło żółw)
    words.sort { |x, y| Example::Crystal.word_cmp(x.downcase, y.downcase) }.should eq(expected)
  end
end
