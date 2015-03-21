#! /usr/bin/env ruby
# -*- coding: UTF-8 -*-

module AI
  class << self

    def next_line
      $stdin.readline.chomp.split(" ").map(&:to_i)
    end

    def write(s)
      $stdout.write s
      $stdout.flush
    end

    def log(s)
      s += "\n" unless s.end_with? "\n"
      @_log.write s
      @_log.flush
    end

    def read_message
      log "reading:"

      t, a, p, s = next_line
      log " turn #{t}, #{a} ants per player (#{p} players), status=#{s}"

      a.times do
        id, x, y, _, _, e, a, b = next_line
        log "  - ant ##{id} (#{x},#{y}) energy=#{e} acid=#{a} brain=#{b}"
      end

      n = next_line.first
      log " #{n} other ants"

      n.times do
        x, y, _, _, b = next_line
        log "  - ant in (#{x}, #{y}), brain=#{b}"
      end

      w, h, n = next_line
      log " map #{w}x#{h}, #{n} points"

      n.times do
        x, y, c, s = next_line
        log "  - (#{x}, #{y}) : #{c} (seen=#{s})"
      end
    end

    def turn
      read_message
      log "go!"

      cmd = @ants.map { |a| "#{a}:forward" }.join(",")
      write "#{cmd}\n"
    rescue EOFError
      log "that's the end."
      @end = true
    end

    def run(log)
      @ants = [0]
      if ARGV.size > 0
        @ants = (0..(ARGV[0].to_i-1)).to_a
      end

      File.open(log, "a") do |f|
        @_log = f
        @end = false

        turn until @end
      end
    end
  end
end

AI.run "/tmp/rb-ai.log"
